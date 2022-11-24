use crate::signalboost::Person;
use maud::{html, Markup, Render};
use serde::{Deserialize, Serialize};
use std::{
    collections::HashMap,
    fmt::{self, Display},
    path::PathBuf,
};

#[derive(Clone, Deserialize, Default)]
pub struct Config {
    pub signalboost: Vec<Person>,
    pub authors: HashMap<String, Author>,
    #[serde(rename = "defaultAuthor")]
    pub default_author: Author,
    pub port: u16,
    #[serde(rename = "clackSet")]
    pub clack_set: Vec<String>,
    #[serde(rename = "resumeFname")]
    pub resume_fname: PathBuf,
    #[serde(rename = "miToken")]
    pub mi_token: String,
    #[serde(rename = "jobHistory")]
    pub job_history: Vec<Job>,
    #[serde(rename = "seriesDescriptions")]
    pub series_descriptions: Vec<SeriesDescription>,
    #[serde(rename = "seriesDescMap")]
    pub series_desc_map: HashMap<String, String>,
    #[serde(rename = "notableProjects")]
    pub notable_projects: Vec<Link>,
    #[serde(rename = "contactLinks")]
    pub contact_links: Vec<Link>,
}

#[derive(Clone, Deserialize, Serialize, Default)]
pub struct Link {
    pub url: String,
    pub title: String,
    pub description: String,
}

impl Render for Link {
    fn render(&self) -> Markup {
        html! {
            span {
                a href=(self.url) {(self.title)}
                @if !self.description.is_empty() {
                    ": "
                    (self.description)
                }
            }
        }
    }
}

#[derive(Clone, Deserialize, Serialize)]
pub enum StockKind {
    Grant,
    Options,
}

impl Default for StockKind {
    fn default() -> Self {
        StockKind::Options
    }
}

fn schema_context() -> String {
    "http://schema.org/".to_string()
}

fn schema_person_type() -> String {
    "Person".to_string()
}

#[derive(Clone, Deserialize, Serialize, Default)]
pub struct Author {
    #[serde(rename = "@context", default = "schema_context")]
    pub context: String,
    #[serde(rename = "@type", default = "schema_person_type")]
    pub schema_type: String,
    pub name: String,
    #[serde(skip_serializing)]
    pub handle: String,
    #[serde(rename = "image", skip_serializing_if = "Option::is_none")]
    pub pic_url: Option<String>,
    #[serde(rename = "inSystem", skip_serializing)]
    pub in_system: bool,
    #[serde(rename = "jobTitle")]
    pub job_title: String,
    #[serde(rename = "sameAs")]
    pub same_as: Vec<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub url: Option<String>,
}

#[derive(Clone, Deserialize, Serialize, Default)]
pub struct SeriesDescription {
    pub name: String,
    pub details: String,
}

impl Render for SeriesDescription {
    fn render(&self) -> Markup {
        html! {
            span {
                a href={"/blog/series/" (self.name)} { (self.name) }
                ": "
                (self.details)
            }
        }
    }
}

#[derive(Clone, Deserialize, Serialize, Default)]
pub struct Stock {
    pub amount: i32,
    #[serde(rename = "cliffYears")]
    pub cliff_years: i32,
    pub kind: StockKind,
    pub liquid: bool,
    #[serde(rename = "vestingYears")]
    pub vesting_years: i32,
}

#[derive(Clone, Deserialize, Serialize, Default)]
pub struct Location {
    pub city: String,
    #[serde(rename = "stateOrProvince")]
    pub state_or_province: String,
    pub country: String,
    pub remote: bool,
}

#[derive(Clone, Deserialize, Serialize, Default)]
pub struct Salary {
    pub amount: i32,
    pub per: String,
    pub currency: String,
    pub stock: Option<Stock>,
}

impl Display for Salary {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}${}/{}", self.currency, self.amount, self.per)
    }
}

impl Render for Salary {
    fn render(&self) -> Markup {
        if self.stock.is_none() {
            return html! { (maud::display(self)) };
        }

        let stock = self.stock.as_ref().unwrap();
        html! {
            details {
                summary {
                    (maud::display(self))
                }

                p{
                    (stock.amount)
                    " "
                    @if stock.liquid {
                        "liquid"
                    }
                    " "
                    @match stock.kind {
                        StockKind::Options => {
                            "options"
                        },
                        StockKind::Grant => {
                            "granted shares"
                        }
                    }
                    ". Vesting for "
                    (stock.vesting_years)
                    " "
                    @if stock.vesting_years == 1 {
                        "year"
                    } @else {
                        "years"
                    }
                    " "
                    " with a cliff of "
                    (stock.cliff_years)
                    " "
                    @if stock.cliff_years == 1 {
                        "year"
                    } @else {
                        "years"
                    }
                    "."
                }
            }
        }
    }
}

#[derive(Clone, Deserialize, Serialize, Default)]
pub struct Job {
    pub company: Company,
    pub title: String,
    #[serde(rename = "startDate")]
    pub start_date: String,
    #[serde(rename = "endDate")]
    pub end_date: Option<String>,
    #[serde(rename = "daysWorked")]
    pub days_worked: Option<i32>,
    #[serde(rename = "daysBetween")]
    pub days_between: Option<i32>,
    pub salary: Salary,
    #[serde(rename = "leaveReason")]
    pub leave_reason: Option<String>,
    pub locations: Vec<Location>,
    pub highlights: Vec<String>,
    #[serde(rename = "hideFromResume")]
    pub hide_from_resume: bool,
}

#[derive(Clone, Deserialize, Serialize, Default)]
pub struct Company {
    pub name: String,
    pub url: Option<String>,
    pub tagline: String,
    pub location: Location,
    pub defunct: bool,
}

impl Render for Job {
    fn render(&self) -> Markup {
        html! {
            tr {
                td { (self.title) }
                td { (self.start_date) }
                td { (self.end_date.as_ref().unwrap_or(&"current".to_string())) }
                td { (if self.days_worked.is_some() { self.days_worked.as_ref().unwrap().to_string() } else { "n/a".to_string() }) }
                td { (self.salary) }
                td { (self.leave_reason.as_ref().unwrap_or(&"n/a".to_string())) }
            }
        }
    }
}
