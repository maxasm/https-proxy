import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

import { useTheme } from '@mui/material/styles';

const NetworkTable = ({set_active_connection, set_dialog_open})=> {
  const theme = useTheme()
  console.log(theme)

  let array = new Array(100).fill(undefined);

  function handle_on_row_click() {
    set_active_connection({"connection": "test"})
    set_dialog_open(true)
  }

  return (
    <TableContainer component={Paper} sx={{marginTop: "10px"}}>
      <Table sx={{minWidth:"650px"}}>
        <TableHead>
          <TableRow sx={{"&.MuiTableRow-root *": {background: theme.palette.primary.main, color: "white"}}}>
            <TableCell align="center">Method</TableCell>
            <TableCell align="center">Protocol</TableCell>
            <TableCell align="center">Server Name</TableCell>
            <TableCell align="center">Response Code</TableCell>
            <TableCell align="center">Content-Type</TableCell>
          </TableRow>    
        </TableHead>
        <TableBody>
          {
            array.map((_, index)=>(
              <TableRow
                sx={{"&.MuiTableRow-root": {cursor: "pointer"}}}
                key={index}
                onClick={()=>handle_on_row_click()}
              >
                <TableCell>GET</TableCell>
                <TableCell align="center">HTTP/2.0</TableCell>
                <TableCell align="center">www.google.com</TableCell>
                <TableCell align="center">200 OK</TableCell>
                <TableCell align="center">text/html</TableCell>
              </TableRow>
            ))
          }
        </TableBody>
      </Table>
    </TableContainer>
  )
}

export default NetworkTable;
