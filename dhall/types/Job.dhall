let Salary = ./Salary.dhall

in  { Type =
        { company : Text
        , title : Text
        , startDate : Text
        , endDate : Optional Text
        , daysWorked : Optional Natural
        , daysBetween : Optional Natural
        , salary : Salary.Type
        , leaveReason : Optional Text
        }
    , default =
      { company = "Unknown"
      , title = "Unknown"
      , startDate = "0000-01-01"
      , endDate = None Text
      , daysWorked = None Natural
      , daysBetween = None Natural
      , salary = Salary::{=}
      , leaveReason = None Text
      }
    }
