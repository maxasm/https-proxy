import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

import { useTheme } from '@mui/material/styles';

const NetworkTable = ({network_info, set_active_connection, set_dialog_open})=> {
  // get the theme
  const theme = useTheme()

  function handle_on_row_click() {
    set_active_connection({"connection": "test"})
    set_dialog_open(true)
  }

  const Rows = ()=> {
    let rows = [];
    network_info.forEach((value,key)=> {
      let url = new URL(value.path)
      let server_name = url.host

      function handle_on_row_click() {
        set_active_connection(key)       
        set_dialog_open(true)
      }

      rows.push(
        <TableRow key={key} onClick={handle_on_row_click}>
          <TableCell align="center">{value.method}</TableCell>
          <TableCell align="center">{value.protocol}</TableCell>
          <TableCell align="center">{server_name}</TableCell>
          <TableCell align="center">{value.responseinfo.headers ? value.responseinfo.headers["Content-Type"] : "Pending ..."}</TableCell>
          <TableCell align="center">{value.responseinfo.status ? value.responseinfo.status : "Pending ... "}</TableCell>
        </TableRow>
      )
    })
    return rows
  }
  
  return (
    <TableContainer component={Paper} sx={{marginTop: "10px"}}>
      <Table sx={{minWidth:"650px"}}>
        <TableHead>
          <TableRow sx={{"&.MuiTableRow-root *": {background: theme.palette.primary.main, color: "white"}}}>
            <TableCell align="center">Method</TableCell>
            <TableCell align="center">Protocol</TableCell>
            <TableCell align="center">Server Name</TableCell>
            <TableCell align="center">Content Type</TableCell>
            <TableCell align="center">Response Code</TableCell>
          </TableRow>    
        </TableHead>
        <TableBody>
          <Rows/>
        </TableBody>
      </Table>
    </TableContainer>
  )
}

export default NetworkTable;
