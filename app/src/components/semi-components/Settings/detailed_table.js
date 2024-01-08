import '../../../CSS/detailed.css';

function Detailed_table(props)
{

    return (
        <table className = "table">
            <tbody>
                <tr>
                    <th>Warehouse Name</th>
                    <th>Shopify Warehouse Name</th>
                    <th>ID</th>
                    <th>Shopify Location ID</th>
                    <th>Created At</th>
                </tr>
                <>
                    {props.table}
                </>
            </tbody>
        </table>
        
    );
};

Detailed_table.defaultProps = 
{
    table: ''
}
export default Detailed_table;