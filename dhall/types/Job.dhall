let Company = ./Company.dhall

let Salary = ./Salary.dhall

let Location = ./Location.dhall

in  { Type =
        { company : Company.Type
        , title : Text
        , contract : Bool
        , startDate : Text
        , endDate : Optional Text
        , daysWorked : Optional Natural
        , daysBetween : Optional Natural
        , salary : Salary.Type
        , leaveReason : Optional Text
        , locations : List Location.Type
        , highlights : List Text
        , hideFromResume : Bool
        }
    , default =
      { company = Company::{=}
      , title = "Unknown"
      , contract = False
      , startDate = "0000-01-01"
      , endDate = None Text
      , daysWorked = None Natural
      , daysBetween = None Natural
      , salary = Salary::{=}
      , leaveReason = None Text
      , locations = [] : List Location.Type
      , highlights = [] : List Text
      , hideFromResume = False
      }
    }
