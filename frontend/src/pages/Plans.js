import React, { useState, useEffect } from 'react';
import Navbar from '../components/Navbar'
import PlanCard from '../components/PlanCard'
import {Grid, Divider, makeStyles, Typography } from '@material-ui/core'
import { FaListUl } from "react-icons/fa"
import { BsGrid3X3GapFill } from "react-icons/bs";
import { DataGrid } from '@material-ui/data-grid';
import axios from '../services/axios'

const useStyles = makeStyles((theme) => ({
    plansCardBox : {
        width: '800px',
        margin: '0 auto',
        maxHeight: '500px'
    },
    switchIcon: {
        height: '15px',
        float : 'right',
        marginTop: 15,
        marginBottom: 15,
        padding: 5,
    },
    tableHead : {
        backgroundColor: "#c5ecec",
    }

}) )


const Plans = () => {
    const classes = useStyles()
    // true = list, false = grid
    const [tab, setTab] = useState(true)
    const [plans, setPlans] = useState([])

    const changeTab = (val) => {
        setTab(val)
    }

    useEffect(() => {
        async function fetchData(){
            const request = await axios.get('/plans/')
            let resp_plans = request.data.data
            setPlans(resp_plans)
            return plans
        }
        fetchData();
    }, []);

    const dataColumns = [
        {
            field : 'id', headerName: 'ID', width:70
        },
        {
            field : 'name', headerName: 'Name', width:100
        },
        {
            field : 'cost', headerName: 'Price', width:100
        },
        {
            field : 'speed', headerName: 'Speed', width:100
        },
        {
            field : 'duration', headerName: 'Duration', width:120
        },
        {
            field : 'notes', headerName: 'Description',width:500
        },

    ]

    const planCards = (
        <React.Fragment>
                 <Grid container >
                    {
                        plans.map((plan, id) => (
                            <Grid key={id} item xs={4} >
                            <PlanCard planName={plan.name} planDescripton={plan.notes} planPrice={plan.cost} planSpeed={plan.speed} planDuration={plan.duration} />
                            </Grid>
                        ))
                    }
                </Grid>
        </React.Fragment>
    )

    const planTable = (
        <React.Fragment>
            <div style={{ height: 600, width: '100%' }}>
                <DataGrid rows={plans} columns={dataColumns} pageSize={10} checkboxSelection />
            </div>
        </React.Fragment>
    )

    return (

        <div>
            <Navbar />
            <div>
                <Typography variant="h4" style={{fontFamily: 'Segoe UI Symbol', padding: '15px'}} >Our Plans <FaListUl className={classes.switchIcon} onClick={ () => changeTab(true) } /> <BsGrid3X3GapFill className={classes.switchIcon} onClick={ () => changeTab(false) } />  </Typography>
            </div>
            <Divider />
            <div className={classes.plansCardBox} >
            {
                tab ? planTable : planCards
            }
            </div>

        </div>
    )
}

export default Plans