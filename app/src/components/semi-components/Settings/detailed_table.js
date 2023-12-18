import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_table(props)
{

    return (
        <table style = {{left: '40%',top: '17px', marginBottom: '0px',fontSize: '13px'}}>
            <tbody>
                <tr>
                    <th>Warehouse</th>
                    <th>Shopify Location</th>
                </tr>
                <>
                    {props.table}
                </>
            </tbody>
        </table>
        
    );
};

Detailed_table.defaultProps = 
{
    table: ''
}
export default Detailed_table;