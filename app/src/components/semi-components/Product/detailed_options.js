import '../../../CSS/detailed.css';

function Detailed_Options(props)
{
    return (
        <tr className = "product_options">
            <td style= {{width: '50%'}}>{props.Option_Value}</td>
            <td style= {{width: '50%'}}>{props.Option_Name}</td>
        </tr>
    );
};

Detailed_Options.defaultProps = 
{
    Option_Name: 'N/A',
    Option_Value: 'N/A',

}
export default Detailed_Options;
