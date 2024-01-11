import {useEffect} from 'react';
import '../../CSS/detailed.css';

function Detailed_User(props)
{
    useEffect(() => 
    {
        let table = document.querySelector(".table");
        let clip = table.querySelectorAll(".button");
        let clip_info = table.querySelectorAll(".user_info");
        //console.log(clip);
        for(let i = 0; i < clip.length; i++)
        {

            clip[i].addEventListener("click", () => 
            {
                
                let copyText = clip_info[i].innerHTML;

                navigator.clipboard.writeText(copyText);
                setTimeout(() =>
                {
                    let message = document.getElementById("message");

                    message.style.display = "block";
                    message.innerHTML = "Copied to Clipboard";
                    message.style.backgroundColor = "rgb(6, 133, 6)";
                    setTimeout(() =>
                    {
                        message.innerHTML = "";
                        message.style.backgroundColor = "transparent";
                        message.style.display = "none";
                    }, 1000);
                }, 1000);
            });
        }
    }, []);

    return (
        <>
            <table className = "table" style ={{marginBottom: '0px'}}>
                <tbody>
                    <tr>
                        <td style= {{width: '25%'}}>Name</td>
                        <td style= {{width: '25%'}} className = 'user_info'>{props.name}</td>
                        <td style= {{width: '25%'}}><button className = 'button' type = 'button'>Copy to Clipboard</button></td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Email</td>
                        <td style= {{width: '25%'}} className = 'user_info'>{props.email}</td>
                        <td style= {{width: '25%'}}><button className = 'button' type = 'button'>Copy to Clipboard</button></td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Api Key</td>
                        <td style= {{width: '25%'}} className = 'user_info'>{props.api_key}</td>
                        <td style= {{width: '25%'}}><button className = 'button' type = 'button'>Copy to Clipboard</button></td>
                    </tr>
                </tbody>
            </table>
        </>
    );
};

export default Detailed_User;