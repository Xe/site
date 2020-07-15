use warp::{http::Response, Rejection, Reply};
use crate::templates::{Html, RenderRucte};

include!(concat!(env!("OUT_DIR"), "/templates.rs"));

pub async fn gitea(pkg_name: &str, git_repo: &str) -> Result<impl Reply, Rejection> {
    Response::builder().html(|o| templates::gitea_html(o, pkg_name, git_repo))
}

pub async fn github(pkg_name: &str, git_repo: &str) -> Result<impl Reply, Rejection> {
    Response::builder().html(|o| templates::github_html(o, pkg_name, git_repo))
}
