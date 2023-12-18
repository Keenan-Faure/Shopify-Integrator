import '../../../CSS/detailed.css';

function Detailed_Quantities(props)
{
    return (
        <tr>
            <td className = "quantity_name" style= {{width: '50%'}}>{props.quantity_name}</td>
            <td className = "quantity_value" style= {{width: '50%'}}>{props.quantity_value}</td>
        </tr>
    );
};

Detailed_Quantities.defaultProps = 
{
    quantity_value: 'N/A',
    quantity_name: 'N/A',

}
export default Detailed_Quantities;
