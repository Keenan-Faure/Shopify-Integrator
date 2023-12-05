import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_Table_View(props)
{
    useEffect(()=> 
    {

        
    }, []);

    return (
        <>
            <tr className="order_totals_line">
                <td className="order_totals_headers">
                    Sub total
                </td>
                <td className="order_totals_middle"></td>
                <td className="order_totals_value">2,300,500</td>
            </tr>
            <tr className="order_totals_line">
                <td className="order_totals_headers">
                    Tax
                </td>
                <td className="order_totals_middle">10%</td>
                <td className="order_totals_value">4,500</td>
            </tr>
            <tr className="order_totals_line">
                <td className="order_totals_headers">
                    Shipping
                </td>
                <td className="order_totals_middle">Standard Shipping</td>
                <td className="order_totals_value">500.00</td>
            </tr>
            <tr className="order_totals_line">
                <td className="order_totals_headers">
                    Total
                </td>
                <td className="order_totals_middle"></td>
                <td className="order_totals_value">2,350,500</td>
            </tr>

        </>
    );
}

export default Detailed_Table_View;