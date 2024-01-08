import {useEffect} from 'react';
import '../../../CSS/page2.css';

function Detailed_Restriction(props)
{
    useEffect(()=> 
    {
    
    }, []);

    return (
        <>
            <tr>
                <td style = {{width: '33%'}} className = "restriction_name">{props.Name}</td>
                <td style = {{width: '33%'}} className = "flag_value" >{props.Flag}</td>
                <td style = {{width: '33%'}}>
                    <label className = "switch" style = {{left: '50%', transform: 'translate(-50%)', marginLeft: '0px'}} id = "slider">
                        <input type = "checkbox" id = "lide" />
                        <span className = "slider round" />
                    </label>
                </td>
            </tr>
        </>
    );
}

export default Detailed_Restriction;