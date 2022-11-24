use crate::app::config::Link;
use serde::Deserialize;

#[derive(Clone, Deserialize)]
pub struct Person {
    pub name: String,
    pub tags: Vec<String>,
    pub links: Vec<Link>,
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
