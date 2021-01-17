use eyre::Result;

#[tokio::main]
async fn main() -> Result<()> {
    kankyo::init()?;

    let key = std::env::var("CF_TOKEN")?;
    let zone_id = std::env::var("CF_ZONE_ID")?;

    let cli = cfcache::Client::new(key, zone_id)?;
    cli.purge(vec!["https://christine.website/.within/health".to_string()])
        .await?;

    Ok(())
}
