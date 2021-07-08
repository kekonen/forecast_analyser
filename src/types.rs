use chrono::{DateTime, Local};

#[derive(Debug)]
pub enum Condition {
    Sunny,
    PartiallyCloudy,
    Cloudy,
    Raining,
    Stormy,
}

#[derive(Debug, Eq, PartialEq)]
pub enum TargetDate {
    Now,
    Hourly(DateTime<Local>),
    Daily(DateTime<Local>),
}

#[derive(Debug)]
pub struct Forecast {
    pub collected_at: DateTime<Local>,
    pub target_at: TargetDate,
    pub temperature: f32,
    pub condition: Condition,
    pub location: String,
    pub source: String,
}

impl Forecast {
    pub fn new(collected_at: DateTime<Local>, target_at: TargetDate, temperature: f32, condition: Condition, location: &str, source: &str) -> Self {
        Self {
            collected_at,
            target_at,
            temperature,
            condition,
            location: location.into(),
            source: source.into(),
        }
    }

    pub fn is_now(&self) -> bool {
        self.target_at == TargetDate::Now
    }
}
