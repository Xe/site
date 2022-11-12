use color_eyre::Result;
use std::{env, fs};
use tracing::debug;
use xesite_types::mastodon::{Toot, User};

#[tokio::main]
async fn main() -> Result<()> {
    color_eyre::install()?;
    tracing_subscriber::fmt::init();

    let args: Vec<String> = env::args().collect();
    debug!("{args:?}");
    if args.len() != 2 {
        eprintln!("Usage: {} <mastodon post URL>", args[0]);
    }

    let mut post_url = args[1].clone();

    let cli = reqwest::Client::builder()
        .user_agent("github.com/Xe/site fetch_mastodon_post")
        .build()?;

    let toot: Toot = cli
        .get(&post_url)
        .header("Accept", "application/json")
        .send()
        .await?
        .error_for_status()?
        .json()
        .await?;

    debug!("got post by {}", toot.attributed_to);

    fs::create_dir_all("./data/toots")?;

    if !post_url.ends_with(".json") {
        post_url = format!("{post_url}.json");
    }
    let post_hash = xesite::hash_string(post_url);

    debug!("wrote post to ./data/toots/{post_hash}.json");

    let mut fout = fs::File::create(&format!("./data/toots/{post_hash}.json"))?;
    serde_json::to_writer_pretty(&mut fout, &toot)?;

    debug!("fetching {} ...", toot.attributed_to);
    let user: User = cli
        .get(&toot.attributed_to)
        .header("Accept", "application/json")
        .send()
        .await?
        .error_for_status()?
        .json()
        .await?;

    fs::create_dir_all("./data/users")?;

    debug!("got user {} ({})", user.preferred_username, user.name);

    let user_url = format!("{}.json", toot.attributed_to);
    let user_hash = xesite::hash_string(user_url);

    debug!("wrote post to ./data/users/{user_hash}.json");
    let mut fout = fs::File::create(&format!("./data/users/{user_hash}.json"))?;
    serde_json::to_writer_pretty(&mut fout, &user)?;

    Ok(())
}
