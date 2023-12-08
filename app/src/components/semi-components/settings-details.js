import {useEffect} from 'react';
import '../../CSS/page2.css';

function Setting_details(props)
{
    useEffect(()=> 
    {
        
    }, []);

    return (
        <>
            <div className = "setting_main">
                <div className = "_title">{props.Key}</div>
                <div className = "setting-details description">{props.Description}</div>
                <div className = "setting-details id" style = {{display: 'none'}}>{props.id}</div>
                <div className = "setting-details value">
                    Currently Set to:<div style ={{fontWeight:'bold'}}className = "_value">{props.Value}</div>
                </div>
                <input type="setting" placeholder = "New Value" className = "_input"></input>
            </div>
        </>
    );
}

Setting_details.defaultProps = 
{
    Key: 'Sub Title of setting',
    DescripDescription1tion: 'Description of product goes here, as well as any additional information',
    Value: 'Value of setting currently in the api',
    id: 'id of the setting'
}

export default Setting_details;