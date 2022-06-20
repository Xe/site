let StockKind = ./StockKind.dhall

in  { Type =
        { kind : StockKind
        , amount : Natural
        , liquid : Bool
        , vestingYears : Natural
        , cliffYears : Natural
        }
    , default =
      { kind = StockKind.Options
      , amount = 0
      , liquid = False
      , vestingYears = 4
      , cliffYears = 1
      }
    }
