import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_shipping(props)
{
    useEffect(()=> 
    {
        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            navigation.style.left = "0%";
            navigation.style.position = "relative";
            navigation.style.width = "100%";
        }
        
    }, []);

    return (
        <>

            <table>
                <tbody>
                    <tr>
                        <th>Customer Firstname</th>
                        <th>Customer Lastname</th>
                        <th>Customer company</th>
                    </tr>
                    <tr>
                        <td>{props.Customer_firstName}</td>
                        <td>{props.Customer_lastName}</td>
                        <td>{props.Customer_company}</td>
                    </tr>
                </tbody>
            </table> 
            <table>
                <tbody>
                    <tr>
                        <th>Shipping Address 1</th>
                    </tr>
                    <tr>
                        <td>{props.Shipping1}</td>
                    </tr>

                    <tr>
                        <th>Shipping Address 2</th>
                    </tr>
                    <tr>
                        <td>{props.Shipping2}</td>
                    </tr>

                    <tr>
                        <th>Shipping Address 3</th>
                    </tr>
                    <tr>
                        <td>{props.Shipping3}</td>
                    </tr>

                    <tr>
                        <th>Shipping Address 4</th>
                    </tr>
                    <tr>
                        <td>{props.Shipping4}</td>
                    </tr>
                    <tr>
                        <th>Shipping Address 5</th>
                    </tr>
                    <tr>
                        <td>{props.Shipping5}</td>
                    </tr>
                </tbody>
            </table> 
        </>
    );
}
Detailed_shipping.defaultProps = 
{
    Shipping1: 'Address1',
    Shipping2: 'Address2',
    Shipping3: 'Address3',
    Shipping4: 'Address4'
}
export default Detailed_shipping;