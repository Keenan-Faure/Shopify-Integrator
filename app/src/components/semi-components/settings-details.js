import {useEffect} from 'react';
import '../../CSS/page2.css';

function Setting_details(props)
{
    useEffect(()=> 
    {
        let _true = document.querySelectorAll(".true");
        let _false = document.querySelectorAll(".false");
        
        for(let i = 0; i < _true.length; i++)
        {
            _false[i].addEventListener("click", () =>
            {
                _false[i].style.backgroundColor = "rgb(148, 11, 11)";
                _false[i].classList.add(".activated");

                _true[i].style.backgroundColor = "white";
                _true[i].style.color = "black";
                _true[i].classList.remove(".activated");

            });

            _true[i].addEventListener("click", () =>
            {
                _true[i].style.backgroundColor = "rgb(19, 99, 8)";
                _true[i].classList.add(".activated");
                _false[i].style.backgroundColor = "white";
                _false[i].style.color = "black";
                _false[i].classList.remove(".activated");
            });
        }
        

    }, []);

    return (
        <>
            <div className = "setting_main">
                <div className = "tittle">{props.Key}</div> 
                <div className = "_title" style={{display: 'none'}}>{props.Title}</div>
                <div className = "setting-details description" style={{backgroundColor: 'transparent'}}>{props.Description}</div>
                <div className = "setting-details id" style = {{display: 'none'}}>{props.id}</div>
                <div className = "setting-details value">
                    Currently Set to:<div style ={{fontWeight:'bold'}}className = "_value">{props.Value}</div>
                </div>
                <div className = "true-false"><button style ={{display: 'inline-block'}} className = "true">True</button> <button style ={{display: 'inline-block'}} className = "false">False</button></div>
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