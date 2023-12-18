import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_warehousing(props)
{

    return (
        <>
            <tr>
                <td className = "warehouse">{props.Warehouse_Location}</td>
                <td className = "fill-able">{props.Shopify_Location}</td>
            </tr>
        </>
        
    );
};

Detailed_warehousing.defaultProps = 
{
    Warehouse_Location: 'N/A', 
    Shopify_Location: 'N/A'

}
export default Detailed_warehousing;