import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContentText from '@mui/material/DialogContentText';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Divider from '@mui/material/Divider';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import TabContext from '@mui/lab/TabContext';
import TabList from '@mui/lab/TabList';
import TabPanel from '@mui/lab/TabPanel';

import {useState} from "react"

const InfoCard = ({selected_connection, request})=> {
  // check if the response has been updated
  if (!request) {
    if (!(selected_connection.responseinfo.status)) {
      return <Typography> Response Pending ... </Typography>
    }
  }

  let url = new URL(selected_connection.url)

  function Headers() {
    let headers = request ? selected_connection.headers : selected_connection.responseinfo.headers
    let headers_array = []
    for (let k in headers){
      headers_array.push(<Typography> {k}: {headers[k]} </Typography>)
    }
    return headers_array
  }

  function Payload() {
    let headers = request ? selected_connection.headers : selected_connection.responseinfo.headers
    let content_type = headers["Content-Type"]
    let content_length = headers["Content-Length"]

    // headers are returned as a type map[string][]string
    function contains(headers, str) {
      let hd
      try {
        hd = headers[0]
      } catch(e) {
        return false
      }
      return hd.indexOf(str) !== -1    
    }

    // helper function to decode base64 content
    function decode_b64() {
      let txt
      try {
        txt = atob(request ? selected_connection.payload : selected_connection.responseinfo.payload) 
        // only show the first 100 characters
        // txt = txt.slice(0, 100)
      } catch(e) {
        txt = "failed to decoded base64 payload"
      }
      return txt
    }
    
    function Preview() {
        return (<Typography> <pre> {decode_b64()} </pre> </Typography>)
    }

    return (
      <Box>
        <Typography variant="h1">Payload</Typography>
        <Divider sx={{marginTop: "10px"}}/>
        <Typography> Content-Type:  {content_type ? content_type : "not stated"} </Typography>
        <Typography> Content-Length:  {content_length ? content_length+" bytes" : "not stated"} </Typography>
        <Divider sx={{marginTop: "10px"}}/>
        <Box>
          <Typography variant="h2"> Payload Preview </Typography>
          <Preview/>
          <Button variant="contained" sx={{marginTop: "20px"}}> Open Using External App </Button>
        </Box>
      </Box>
    )
  }

  return (
    <Box>
      <Typography variant="h1"> {request ? "HTTP Request" : "HTTP Response"} </Typography>
      {request && <Divider sx={{marginTop: "10px"}}/>}
      {request && <Typography>{selected_connection.method} {selected_connection.protocol}</Typography> }
      {request && <Typography>{selected_connection.path}</Typography> }
      {!request && <Typography>{selected_connection.responseinfo.status}</Typography> }
      <Typography variant="h1"> Headers </Typography>
      <Box>
        <Headers/>
      </Box>
      <Box>
        <Payload/>
      </Box>
    </Box>
  )
}

// the network infor comes from the server via websockets
// it contains all information about the request and response
const NetworkInspect = ({dialog_open, set_dialog_open, network_info, active_connection, set_active_connection})=> {

  const [active_tab, set_active_tab] = useState("request")
  function handle_on_tab_change() {
    set_active_tab(active_tab === "request" ? "response" : "request")
  }
  
  function handle_on_dialog_close() {
    set_dialog_open(false)  
    // set_active_connection("")
  }
  
  let selected_connection = network_info.get(active_connection)
  
  return (
  <Dialog
    fullWidth
    maxWidth="1000px"
    open={dialog_open}
    onClose={handle_on_dialog_close}>
      <DialogTitle>
        <Typography variant="h1">Connection Information</Typography>
      </DialogTitle>
      <DialogContent>
      <TabContext value={active_tab}>
        <Box>
          <TabList onChange={handle_on_tab_change}>
            <Tab label="request" value="request"/>
            <Tab label="response" value="response"/>
          </TabList>
        </Box>
        <TabPanel value="request">
          <InfoCard selected_connection={selected_connection} request={true}/>
        </TabPanel>
        <TabPanel value="response">
          <InfoCard selected_connection={selected_connection} request={false}/>
        </TabPanel>
      </TabContext>
      </DialogContent>
  </Dialog>
  )
}

export default NetworkInspect;
