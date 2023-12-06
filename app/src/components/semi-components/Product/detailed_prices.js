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
                        <td>{props.Price_Name}</td>
                        <td>{props.Price_Value}</td>
                    </tr>
                </tbody>
            </table>
        </>
    );
}

Detailed_Price.defaultProps = 
{
    Price_Name: 'Name',
    Price_Value: 'Value'
}
export default Detailed_Price;
