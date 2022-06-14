use serde::Deserialize;

#[derive(Clone, Debug, Deserialize)]
pub struct Person {
    pub name: String,
    pub tags: Vec<String>,

    #[serde(rename = "gitLink")]
    pub git_link: Option<String>,

    pub twitter: Option<String>,
}

#[cfg(test)]
mod tests {
    use color_eyre::eyre::Result;
    #[test]
    fn load() -> Result<()> {
        let _people: Vec<super::Person> =
            serde_dhall::from_file("./dhall/signalboost.dhall").parse()?;

        Ok(())
    }
}
