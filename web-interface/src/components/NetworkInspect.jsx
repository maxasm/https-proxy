import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContentText from '@mui/material/DialogContentText';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Box from '@mui/material/Box';
import TabContext from '@mui/lab/TabContext';
import TabList from '@mui/lab/TabList';
import TabPanel from '@mui/lab/TabPanel';

import {useState} from "react"

// the network infor comes from the server via websockets
// it contains all information about the request and response
const NetworkInspect = ({dialog_open, set_dialog_open, network_info})=> {

  const [active_tab, set_active_tab] = useState("request")
  function handle_on_tab_change() {
    set_active_tab(active_tab === "request" ? "response" : "request")
  }
  
  return (
  <Dialog
    fullWidth
    maxWidth="1000px"
    open={dialog_open}
    onClose={()=>set_dialog_open(false)}>
      <DialogTitle>Connection Information</DialogTitle>
      <DialogContent>
      <TabContext value={active_tab}>
        <Box>
          <TabList onChange={handle_on_tab_change}>
            <Tab label="request" value="request"/>
            <Tab label="response" value="response"/>
          </TabList>
        </Box>
        <TabPanel value="request">Request Information</TabPanel>
        <TabPanel value="response">Response Information</TabPanel>
      </TabContext>
      </DialogContent>
  </Dialog>
  )
}

export default NetworkInspect;
