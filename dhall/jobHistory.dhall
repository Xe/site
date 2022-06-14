let Job = ./types/Job.dhall

let Salary = ./types/Salary.dhall

let annual = \(rate : Natural) -> Salary::{ amount = rate }

let hourly = \(rate : Natural) -> Salary::{ amount = rate, per = "hour" }

let annualCAD = \(rate : Natural) -> Salary::{ amount = rate, currency = "CAD" }

in  [ Job::{
      , company = "Symplicity"
      , title = "Junior Systems Administrator"
      , startDate = "2013-11-11"
      , endDate = Some "2014-01-06"
      , daysWorked = Some 56
      , salary = annual 50000
      , leaveReason = Some "terminated"
      }
    , Job::{
      , company = "OpDemand"
      , title = "Software Engineering Intern"
      , startDate = "2014-07-14"
      , endDate = Some "2014-08-27"
      , daysWorked = Some 44
      , daysBetween = Some 189
      , salary = annual 35000
      , leaveReason = Some "terminated"
      }
    , Job::{
      , company = "Crowdflower (contract)"
      , title = "Consultant"
      , startDate = "2014-09-17"
      , endDate = Some "2014-10-15"
      , daysWorked = Some 28
      , daysBetween = Some 21
      , salary = hourly 90
      , leaveReason = Some "contract not renewed"
      }
    , Job::{
      , company = "VTCSecure (contract)"
      , title = "Consultant"
      , startDate = "2014-10-27"
      , endDate = Some "2015-02-09"
      , daysWorked = Some 105
      , daysBetween = Some 12
      , salary = hourly 90
      , leaveReason = Some "contract not renewed"
      }
    , Job::{
      , company = "IMVU"
      , title = "Site Reliability Engineer"
      , startDate = "2015-03-30"
      , endDate = Some "2016-03-07"
      , daysWorked = Some 343
      , daysBetween = Some 49
      , salary = annual 125000
      , leaveReason = Some "demoted"
      }
    , Job::{
      , company = "IMVU"
      , title = "Systems Administrator"
      , startDate = "2016-03-08"
      , endDate = Some "2016-04-01"
      , daysWorked = Some 24
      , daysBetween = Some 1
      , salary = annual 105000
      , leaveReason = Some "quit"
      }
    , Job::{
      , company = "Pure Storage"
      , title = "Member of Technical Staff"
      , startDate = "2016-04-04"
      , endDate = Some "2016-08-03"
      , daysWorked = Some 121
      , daysBetween = Some 3
      , salary = annual 135000
      , leaveReason = Some "quit"
      }
    , Job::{
      , company = "Backplane.io (defunct)"
      , title = "Software Engineer"
      , startDate = "2016-08-24"
      , endDate = Some "2016-11-22"
      , daysWorked = Some 90
      , daysBetween = Some 21
      , salary = annual 105000
      , leaveReason = Some "terminated"
      }
    , Job::{
      , company = "Heroku (contract)"
      , title = "Consultant"
      , startDate = "2017-02-13"
      , endDate = Some "2017-11-13"
      , daysWorked = Some 273
      , daysBetween = Some 83
      , salary = hourly 120
      , leaveReason = Some "hired"
      }
    , Job::{
      , company = "Heroku"
      , title = "Senior Software Engineer"
      , startDate = "2017-11-13"
      , endDate = Some "2019-03-08"
      , daysWorked = Some 480
      , daysBetween = Some 0
      , salary = annual 150000
      , leaveReason = Some "quit"
      }
    , Job::{
      , company = "Lightspeed POS"
      , title = "Expert principal en fiabilité du site"
      , startDate = "2019-05-06"
      , endDate = Some "2020-11-27"
      , daysWorked = Some 540
      , daysBetween = Some 48
      , salary = annualCAD 115000
      , leaveReason = Some "quit"
      }
    , Job::{
      , company = "Tailscale"
      , title = "Software Designer"
      , startDate = "2020-12-14"
      , endDate = Some "2022-03-01"
      , daysWorked = Some 442
      , daysBetween = Some 0
      , salary = annualCAD 135000
      , leaveReason = Some "raise"
      }
    , Job::{
      , company = "Tailscale"
      , title = "Archmage of Infrastructure"
      , startDate = "2022-03-01"
      , salary = annualCAD 147150
      }
    ]
