use crate::app::Config;
use maud::{html, Markup};
use std::sync::Arc;

pub mod nag;

pub fn salary_history(cfg: Arc<Config>) -> Markup {
    html! {
        table.salary_history {
            tr {
                th { "Title" }
                th { "Start Date" }
                th { "End Date" }
                th { "Days Worked" }
                th { "Salary" }
                th { "How I Left" }
            }
            @for job in &cfg.clone().job_history {
                (job.pay_history_row())
            }
        }
    }
}
