import '../../../CSS/detailed.css';

function Detailed_Warehousing(props)
{
    return (
        <div>
            <div className = "Warehose-title">{props.Warehouse_Title}</div>
            <div className = "Warehouse-Text">{props.Warehouse_Text}</div>
            
        </div>
    );
};

Detailed_Warehousing.defaultProps = 
{
    Warehouse_Text: 'Text for warehouse goes here',
    Warehouse_Title: 'Title for warehouse goes here'
}
export default Detailed_Warehousing;
