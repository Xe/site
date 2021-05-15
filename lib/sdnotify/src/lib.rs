//! Notify service manager about start-up completion and
//! other daemon status changes.
//!
//! ### Prerequisites
//!
//! A unit file with service type `Notify` is required.
//!
//! Example:
//! ```toml
//! [Unit]
//! Description=Frobulator
//! [Service]
//! Type=notify
//! ExecStart=/usr/sbin/frobulator
//! [Install]
//! WantedBy=multi-user.target
//! ```
//! ### Sync API
//! ```no_run
//!     use sdnotify::{SdNotify, Message, Error};
//!
//! # fn notify() -> Result<(), Error> {
//!     let notifier = SdNotify::from_env()?;
//!     notifier.notify_ready()?;
//! #   Ok(())
//! # }
//! ```
//!
//! ### Async API
//! ```no_run
//!     use sdnotify::{Message, Error, async_io::SdNotify};
//!     use tokio::prelude::*;
//!     use tokio::runtime::current_thread::Runtime;
//!
//! # fn notify() -> Result<(), Error> {
//!     let notifier = SdNotify::from_env()?;
//!     let mut rt = Runtime::new().unwrap();
//!     rt.block_on(notifier.send(Message::ready())).unwrap();
//! #   Ok(())
//! # }
//! ```

use std::env;
use std::os::unix::net::UnixDatagram;
use std::path::Path;

/// Message to send to init system
#[derive(Debug)]
pub struct Message(InnerMessage);

impl Message {
    /// Tells the init system that daemon startup is finished.
    pub fn ready() -> Self {
        Message(InnerMessage::Ready)
    }

    /// Passes a single-line status string back to the init system that describes the daemon state.
    pub fn status(status: String) -> Result<Self, std::io::Error> {
        if status.as_bytes().iter().any(|x| *x == b'\n') {
            return Err(std::io::Error::new(
                std::io::ErrorKind::InvalidData,
                "newline not allowed",
            ));
        }
        Ok(Message(InnerMessage::Status(status)))
    }

    /// Tells systemd to update the watchdog timestamp.
    /// This is the keep-alive ping that services need to issue in regular
    /// intervals if WatchdogSec= is enabled for it.
    pub fn watchdog() -> Self {
        Message(InnerMessage::Watchdog)
    }

    /// Tells systemd what the main pid of this service is.
    /// This is needed in order to hack up 0-downtime deployments.
    pub fn main_pid(pid: u32) -> Self {
        Message(InnerMessage::MainPid(pid))
    }
}

#[derive(Debug)]
enum InnerMessage {
    Ready,
    Status(String),
    Watchdog,
    MainPid(u32),
}

#[derive(Debug)]
pub enum Error {
    NoSocket,
    Io(std::io::Error),
}

impl std::fmt::Display for Error {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Error::NoSocket => write!(f, "NOTIFY_SOCKET variable not set"),
            Error::Io(err) => write!(f, "{}", err),
        }
    }
}

impl From<env::VarError> for Error {
    fn from(_: env::VarError) -> Error {
        Error::NoSocket
    }
}

impl From<std::io::Error> for Error {
    fn from(err: std::io::Error) -> Error {
        Error::Io(err)
    }
}

impl std::error::Error for Error {}

pub struct SdNotify(UnixDatagram);

impl SdNotify {
    pub fn from_env() -> Result<Self, Error> {
        let sockname = env::var("NOTIFY_SOCKET")?;
        Self::from_path(sockname)
    }

    pub fn from_path<P: AsRef<Path>>(path: P) -> Result<Self, Error> {
        let socket = UnixDatagram::unbound()?;
        socket.connect(path)?;
        Ok(SdNotify(socket))
    }

    /// Tells the init system that daemon startup is finished.
    pub fn notify_ready(&self) -> Result<(), std::io::Error> {
        self.state(Message::ready())
    }

    /// Passes a single-line status string back to the init system that describes the daemon state.
    pub fn set_status(&self, status: String) -> Result<(), std::io::Error> {
        self.state(Message::status(status)?)
    }

    /// Tells systemd to update the watchdog timestamp.
    /// This is the keep-alive ping that services need to issue in regular
    /// intervals if WatchdogSec= is enabled for it.
    pub fn ping_watchdog(&self) -> Result<(), std::io::Error> {
        self.state(Message::watchdog())
    }

    /// Tells systemd what the main pid of this service is.
    /// This is needed in order to hack up 0-downtime deployments.
    pub fn set_main_pid(&self, pid: u32) -> Result<(), std::io::Error> {
        self.state(Message::main_pid(pid))
    }

    fn state(&self, state: Message) -> Result<(), std::io::Error> {
        match state.0 {
            InnerMessage::Ready => self.0.send(b"READY=1")?,
            InnerMessage::Status(status) => self.0.send(format!("STATUS={}", status).as_bytes())?,
            InnerMessage::Watchdog => self.0.send(b"WATCHDOG=1")?,
            InnerMessage::MainPid(pid) => self.0.send(format!("MAINPID={}", pid).as_bytes())?,
        };
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    #[test]
    fn ok() {
        use super::*;

        let path = "/tmp/kek-async.sock";

        let _ = std::fs::remove_file(path);

        let listener = UnixDatagram::bind(path).unwrap();
        let notifier = SdNotify::from_path(path).unwrap();
        notifier.state(Message::ready()).unwrap();
        let mut buf = [0; 100];
        listener.recv(&mut buf).unwrap();
        assert_eq!(&buf[..7], b"READY=1");
    }
}
