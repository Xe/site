use log::{debug, error, info, trace, warn};

fn main() {
    pretty_env_logger::init();

    trace!("starting main");
    debug!("debug message");
    info!("this is some information");
    warn!("oh no something bad is about to happen");
    error!("oh no it's an error");
}
