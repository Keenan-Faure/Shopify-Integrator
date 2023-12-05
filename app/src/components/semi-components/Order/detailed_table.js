import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_Table(props)
{
    useEffect(()=> 
    {

        
    }, []);

    return (
        <>
            <tr className="order_line">
                <td className="order_line_image"></td>
                <td className="order_line_title">
                    <a>{props.Order_Title}</a>
                    <p><b>{props.Order_SKU}</b></p>
                </td>
                <td className="order_line_qty">x{props.Quantity}</td>
                <td className="order_line_qty">#{props.Barcode}</td>
                <td className="order_line_price">R{props.Order_Price}</td>
            </tr>
        </>
    );
}

export default Detailed_Table;