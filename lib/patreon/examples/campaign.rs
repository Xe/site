use patreon::*;

#[tokio::main]
async fn main() -> Result<()> {
    pretty_env_logger::init();
    let creds: Credentials = envy::prefixed("PATREON_").from_env().unwrap();
    let cli = Client::new(creds);

    let camp = cli.campaign().await?;
    println!("{:#?}", camp);

    let id = camp.data[0].id.clone();

    let pledges = cli.pledges(id).await?;
    println!("{:#?}", pledges);
    Ok(())
}
