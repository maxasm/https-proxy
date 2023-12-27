import AppBar from "@mui/material/AppBar"
import FormControlLabel from "@mui/material/FormControlLabel"
import Paper from "@mui/material/Paper"
import Switch from "@mui/material/Switch"
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'

/** state **/
import {useState} from "react"

const FilterBar = ()=> {
  const [apply_filter, update_apply_filter] = useState(false)

  function handleChange() {
    update_apply_filter(!apply_filter)
  }
  
  return (
    <Paper sx={{display: "flex", flexDirection: "column", padding: "10px", alignItems: "flex-start"}}>
      <FormControlLabel
        control={<Switch checked={apply_filter} onChange={handleChange}/>}
        label={apply_filter ? "Disable Filter" : "Enable Filter"}
        sx={{"display": "block"}}
      >
      </FormControlLabel>
      <TextField
        multiline
        label="Filter"
        sx={{alignSelf: "stretch"}}
        inputProps={{ spellCheck: false }}
        error={true}
        helperText={"Handle possible errors"}
      />
      <Button
        variant="contained"
        sx={{marginTop: "6px"}}
        disabled={!apply_filter}
      > Apply Filter </Button>
    </Paper>
  )
}

export default FilterBar;
