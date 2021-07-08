use crate::{TargetDate, Forecast, Condition};
use std::error;
use chrono::{DateTime, Local};

pub trait ForecastFetcher {
    fn fetch(&self, cities: Vec<String>) -> Result<Vec<Forecast>, Box<dyn error::Error>>;
    fn get_available_cities(&self, ) -> Vec<String>;
    // fn prepare_for_cities(&self);
}

pub struct WttrInFetcher {
}

impl ForecastFetcher for WttrInFetcher {
    fn fetch(&self, cities: Vec<String>) -> Result<Vec<Forecast>, Box<dyn error::Error>> {
        let available_cities = self.get_available_cities();
        let mut going_to_do_cities = cities.clone();
        going_to_do_cities.retain(|city|available_cities.contains(city));
        Ok(vec![
            Forecast::new(Local::now(), TargetDate::Now, 20.0, Condition::Cloudy, "Berlin", "accuweather")
        ])
    }
    fn get_available_cities(&self, ) -> Vec<String> {
        vec![
            "Berlin".into()
            ]
    }
}