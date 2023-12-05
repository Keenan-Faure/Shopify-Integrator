import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_Address(props)
{
    useEffect(()=> 
    {

        
    }, []);

    return (
        <>
            <div className="order_data_customer_address_view">
                <div className="customer_address_title">{props.Address_Name}</div>
                <hr />
                <div className="customer_address_header">Address Line 1</div>
                <div className="customer_address_tiles">{props.Address1}</div>
                <div className="customer_address_header">AddressLine 2</div>
                <div className="customer_address_tiles">{props.Address2}</div>
                <div className="customer_address_header">City</div>
                <div className="customer_address_tiles">{props.Address3}</div>
                <div className="customer_address_header">Suburb</div>
                <div className="customer_address_tiles">{props.Address4}</div>
                <div className="customer_address_header">Postal Code</div>
                <div className="customer_address_tiles">{props.Address5}</div>
            </div>
        </>
    );
}
Detailed_Address.defaultProps = 
{
    Address_Name: 'Address Name', 
    Address1: 'Address Line1',
    Address2: 'Address Line2',
    Address3: 'Address Line3',
    Address4: 'Address Line4',
    Address5: 'Address Line5'
}
export default Detailed_Address;
