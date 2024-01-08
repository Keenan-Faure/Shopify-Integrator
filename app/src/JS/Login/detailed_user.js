import '../../CSS/detailed.css';

function Detailed_User(props)
{
    return (
        <>
            <table className = "table" style ={{marginBottom: '0px'}}>
                <tbody>
                    <tr>
                        <th style= {{width: '25%'}}>Id</th>
                        <td style= {{width: '25%'}}>{props.id}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Name</td>
                        <td style= {{width: '25%'}}>{props.name}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Email</td>
                        <td style= {{width: '25%'}}>{props.email}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Api_Key</td>
                        <td style= {{width: '25%'}}>{props.api_key}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Webhook_Token</td>
                        <td style= {{width: '25%'}}>{props.web_token}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Created At:</td>
                        <td style= {{width: '25%'}}>{props.created}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Updated At:</td>
                        <td style= {{width: '25%'}}>{props.updated}</td>
                    </tr>
                </tbody>
            </table>
        </>
    );
};

export default Detailed_User;