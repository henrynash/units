# 

`js
let units = window.Antha.units;

{
  let r = units.Parse(1, "ml");

  console.log(1.0 === r[0].Quantity()
    && "ml" === r[0].MeasurementUnit()
    && null === r[1]);
}

{
  let r = units.Parse(1, "badunit");

  console.log(0.0 === r[0].Quantity()
    && "" === r[0].MeasurementUnit()
    && "" !== r[1].Error());
}

{
  let r = units.New("ml", units.Measurement(1, "L"));

  console.log(1000.0 === r[0].Quantity() 
    && "ml" === r[0].MeasurementUnit()
    && null === r[1]);
}

{
  let r = units.New("g", units.Measurement(1, "L"), units.Measurement(2, "g/L"));

  console.log(2.0 === r[0].Quantity() 
    && "g" === r[0].MeasurementUnit()
    && null === r[1]);
}
```
