import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_Table_View()
{
    useEffect(()=> 
    {

        
    }, []);

    return (
        <>
            <tr className="order_totals_line">
                <td className="order_totals_headers">
                    {props.Total_Heading}
                </td>
                <td className="order_totals_middle">{props.Total_Middle}</td>
                <td className="order_totals_value">{props.subTotal}</td>
            </tr>
        </>
    );
}

export default Detailed_Table_View;