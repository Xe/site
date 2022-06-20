let Stock = ./Stock.dhall

in  { Type =
        { amount : Natural
        , currency : Text
        , per : Text
        , stock : Optional Stock.Type
        }
    , default =
      { amount = 0, currency = "USD", per = "year", stock = None Stock.Type }
    }
