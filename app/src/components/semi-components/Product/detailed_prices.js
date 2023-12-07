import '../../../CSS/detailed.css';

function Detailed_Price(props)
{
    return (
        <>
            <table>
                <tbody>
                    <tr>
                        <th>Name</th>
                        <th>Price</th>
                    </tr>
                    <tr>
                        <td className = "price_name">{props.Price_Name}</td>
                        <td className = "price_value">{props.Price_Value}</td>
                    </tr>
                </tbody>
            </table>
        </>
    );
}

Detailed_Price.defaultProps = 
{
    Price_Name: 'N/A',
    Price_Value: 'N/A'
}
export default Detailed_Price;
