import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_warehousing(props)
{

    return (
        <>
            <tr className = "warehouse">
                <td>{props.Warehouse_Name}</td>
                <td>{props.Shopify_Warehouse_Name}</td>
                <td>{props.id}</td>
                <td>{props.Shopify_Location_ID}</td>
                <td>{props.Created_At}</td>
            </tr>
        </>
        
    );
};

Detailed_warehousing.defaultProps = 
{
    Warehouse_Name: 'N/A', 
    Shopify_Warehouse_Name: 'N/A',
    id: 'N/A',
    Shopify_Location_ID: 'N/A',
    Created_At: 'N/A'

}
export default Detailed_warehousing;