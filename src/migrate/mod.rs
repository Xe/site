use super::establish_connection;
use color_eyre::eyre::Result;
use rusqlite_migration::{Migrations, M};

#[instrument(err)]
pub fn run() -> Result<()> {
    info!("running");
    let mut conn = establish_connection()?;

    let migrations = Migrations::new(vec![M::up(include_str!("./base_schema.sql"))]);
    conn.pragma_update(None, "journal_mode", &"WAL").unwrap();

    migrations.to_latest(&mut conn)?;

    Ok(())
}
