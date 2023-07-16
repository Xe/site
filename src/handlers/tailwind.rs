use std::sync::Arc;
use axum::response::IntoResponse;
use axum::{Extension, extract::Path};
use encre_css::Config;
use http::{HeaderMap, header};
use crate::app::State;
use super::blog;

#[instrument(skip(state, headers))]
pub async fn post_view(
    Path(name): Path<String>,
    Extension(state): Extension<Arc<State>>,
    headers: HeaderMap,
) -> impl IntoResponse {
    let mut config = Config::default();

    encre_css_typography::register(&mut config);
    
    let (_code, body) = blog::post_view(Path(name), Extension(state), headers).await.unwrap();

    ([(header::CONTENT_TYPE, "text/css")], encre_css::generate([body.0.as_str()], &config))
}
