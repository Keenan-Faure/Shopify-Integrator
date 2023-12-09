import '../../../CSS/detailed.css';

function Detailed_Options(props)
{
    return (
            <>
                    <th>{props.Option_Name}</th>
                    <td>{props.Option_Value}</td>
            </>

           

            
    );
};

Detailed_Options.defaultProps = 
{
    Option_Name: 'N/A',
    Option_Value: 'N/A',

}
export default Detailed_Options;
