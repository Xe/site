use crate::templates::RenderRucte;
use warp::{http::Response, Rejection, Reply};

include!(concat!(env!("OUT_DIR"), "/templates.rs"));

pub async fn gitea(pkg_name: &str, git_repo: &str, branch: &str) -> Result<impl Reply, Rejection> {
    Response::builder().html(|o| templates::gitea_html(o, pkg_name, git_repo, branch))
}

pub async fn github(pkg_name: &str, git_repo: &str, branch: &str) -> Result<impl Reply, Rejection> {
    Response::builder().html(|o| templates::github_html(o, pkg_name, git_repo, branch))
}
