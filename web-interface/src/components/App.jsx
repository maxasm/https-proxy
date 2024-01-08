import Box from "@mui/material/Box"
import Button from "@mui/material/Button"
import "../styles/App.css"

import FilterBar from "./FilterBar.jsx"
import NetworkTable from "./NetworkTable.jsx"
import NetworkInspect from "./NetworkInspect.jsx"

import {useState, useEffect} from "react"

import { ThemeProvider, createTheme } from '@mui/material/styles';

const theme = createTheme({
  typography: {
    fontFamily: "Fantasque Sans Mono",
    h1: {
      fontSize: "32px",    
    },
    h2: {
      fontSize: "24px",
    }
  },
});

const App = ()=> {
  const [dialog_open, set_dialog_open] = useState(false)
  // a map of ids to the info
  const [network_info, set_network_info] = useState(new Map())
  // active connection is the ID of the connection
  const [active_connection, set_active_connection] = useState("")

  // handle updating and setting the values for Network Info
  // via WebSockets.
  useEffect(()=>{
    console.log("running useEffect()")
    async function fetch_network_info() {
      let websocket = new WebSocket(`ws://${location.host}/ws`)

      websocket.onopen = function() {
        console.log("websocket connection open")
      }

      websocket.onclose = function() {
        console.log("websocket connection closed")
      }

      websocket.onmessage = function(event) {
        let data = event.data
        let json_data = JSON.parse(data)
        console.log(json_data)
        set_network_info(()=> {
          let id = json_data.id
          let updated_map = network_info.set(id, json_data)
          return new Map(updated_map)
        })
      }

      websocket.onerror = function(error) {
        console.log(`websocket error: ${error}`)
      }
    }

    fetch_network_info()
  },[]);

  return (
    <ThemeProvider theme={theme}>
      <Box sx={{"padding":"10px"}}>
        <FilterBar
          set_network_info={set_network_info}
        />
        <NetworkTable
          network_info={network_info}
          set_active_connection={set_active_connection}
          set_dialog_open={set_dialog_open}
        />
        <NetworkInspect
          network_info={network_info}
          dialog_open={dialog_open}
          set_dialog_open={set_dialog_open}
          active_connection={active_connection}
          set_active_connection={set_active_connection}
        />
      </Box>
    </ThemeProvider>
  )
}

export default App;
