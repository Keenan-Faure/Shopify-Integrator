import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_Table_View(props)
{
    useEffect(()=> 
    {

        
    }, []);

    return (
        <>
            <tr>
                <td >
                    {props.Total_Heading}
                </td>
                <td>{props.Total_Middle}</td>
                <td>{props.subTotal}</td>
            </tr>
        </>
    );
}

export default Detailed_Table_View;