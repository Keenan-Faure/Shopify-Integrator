import {useEffect} from 'react';
import '../../CSS/page1.css';

function Detailed_customer(props)
{
    useEffect(()=> 
    {
        
    }, []);

    return (

        <div className = "pan">
            
        </div>
    );
};

Detailed_customer.defaultProps = 
{
    Customer_Title: 'Customer title',
    Customer_Code: 'Customer code',
    Customer_Options: 'Options',
    Customert_Category: 'Category',
    Customer_Type: 'Type',
    Customer_Vendor: 'Vendor',
    Customer_Price: 'Price'
}
export default Detailed_customer;