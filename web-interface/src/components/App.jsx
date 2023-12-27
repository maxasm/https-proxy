import Box from "@mui/material/Box"
import Button from "@mui/material/Button"
import "../styles/App.css"

import FilterBar from "./FilterBar.jsx"
import NetworkTable from "./NetworkTable.jsx"
import NetworkInspect from "./NetworkInspect.jsx"

import {useState, useEffect} from "react"

import { ThemeProvider, createTheme } from '@mui/material/styles';

const theme = createTheme({
  typography: {fontFamily: "Fantasque Sans Mono"},
})

const App = ()=> {
  const [dialog_open, set_dialog_open] = useState(false)
  const [network_info, set_network_info] = useState({})
  const [active_connection, set_active_connection] = useState({})

  // handle updating and setting the values for Network Info
  // via WebSockets.
  useEffect(()=>{
    async function fetch_network_info() {
      let websocket = new WebSocket(`ws://${location.host}/ws`)

      websocket.onopen = function() {
        console.log("websocket connection open")

        websocket.send("Hello from Chrome.")
      }

      websocket.onclose = function() {
        console.log("websocket connection closed")
      }

      websocket.onmessage = function(event) {
        let data = event.data
        let json_data = JSON.parse(data)
        console.log(json_data)
      }

      websocket.onerror = function(error) {
        console.log(`websocket error: ${error}`)
      }
    }

    fetch_network_info()
  });

  // TODO: Network info should have a `visible` field which should be toggled
  // to make it filter the rows from the table without removing the rows
  // completely from the array. 

  return (
    <ThemeProvider theme={theme}>
      <Box sx={{"padding":"10px"}}>
        <FilterBar
          set_network_info={set_network_info}
        />
        <NetworkTable
          network_info={{}}
          set_active_connection={set_active_connection}
          set_dialog_open={set_dialog_open}
        />
        <NetworkInspect
          network_info={{}}
          dialog_open={dialog_open}
          set_dialog_open={set_dialog_open}
          active_connection={{}}
          set_active_connection={{}}
        />
      </Box>
    </ThemeProvider>
  )
}

export default App;
