import {useEffect} from 'react';
import '../../CSS/page1.css';

function Detailed_order(props)
{
    useEffect(()=> 
    {
        
    }, []);

    return (

        <div className = "pan">
            
        </div>
    );
};

Detailed_order.defaultProps = 
{
    Order_Title: 'Order title',
    Order_Code: 'Order code',
    Order_Options: 'Options',
    Order_Category: 'Category',
    Order_Type: 'Type',
    Order_Vendor: 'Vendor',
    Order_Price: 'Price'
}
export default Detailed_order;