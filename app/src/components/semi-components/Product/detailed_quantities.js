import '../../../CSS/detailed.css';

function Detailed_Quantities(props)
{
    return (
        <div>
            <div className = "quantity_value">{props.quantity_value}</div>
            <div className = "quantity_name">{props.quantity_name}</div>
        </div>
    );
};

Detailed_Quantities.defaultProps = 
{
    quantity_value: 'N/A',
    quantity_name: 'N/A',

}
export default Detailed_Quantities;
