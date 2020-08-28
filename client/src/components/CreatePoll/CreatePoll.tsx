import React from 'react'
import { Formik, Form, Field, FieldArray, ErrorMessage } from 'formik'
import styles from './CreatePoll.module.css'
import * as Yup from 'yup'
import { IPollToSend } from '../../models/poll'
import { ApiCreatePoll } from '../../api/api'
import { RouteComponentProps } from 'react-router'



class CreatePoll extends React.Component<RouteComponentProps, any> {
    render = () => {
        let initialValues: IPollToSend = {
            choices: new Array(2).fill(""),
            name: "",
            settings: {
                allowMultiple: 1,
                timeoutMinutes: 10,
                filter: "",
            },
        }

        const validationSchema = Yup.object().shape({
            name: Yup.string().required().max(100),
            choices: Yup.array().of(
                Yup.string().required()
            ).min(2).test('unique', 'must be unique', values => {
                return new Set(values).size === values?.length
            }),
            settings: Yup.object().shape({
                allowMultiple: Yup.number().min(1),
                timeoutMinutes: Yup.number().min(1).max(120),
                filter: Yup.string().oneOf(["", "ip", "cookie"]),
            })
        })

        return (
            <div className="block">
                <h2>Create Poll</h2>
                <Formik initialValues={initialValues} onSubmit={values => {
                    console.log(`posting`, values)
                    // hack, field with type="number stores value as string"
                    if (typeof values.settings.allowMultiple === 'string') values.settings.allowMultiple = Number.parseInt(values.settings.allowMultiple)
                    ApiCreatePoll(values).then((res) => {
                        this.props.history.push(`/poll/${res.data.uuid}`)
                    }).catch(err => { alert("couldn't create poll: " + err) })
                }} render={
                    ({ values }) => {
                        return (
                            <Form>
                                <div>Poll name</div>
                                <Field name="name"></Field>
                                <ErrorMessage name="name" render={message => (
                                    <div className={styles.error}>{message}</div>
                                )}></ErrorMessage>
                                <div>Choices</div>
                                <FieldArray name="choices">{
                                    ({ remove, push }) => (
                                        <div>
                                            {values.choices.map((_choice, idx, choices) => {
                                                return (
                                                    <div>
                                                        <Field name={`choices.${idx}`}></Field>
                                                        {choices.length > 2 && <button onClick={(event) => {
                                                            event.preventDefault()
                                                            remove(idx)
                                                            values.settings.allowMultiple--
                                                        }}>-</button>}
                                                        {<ErrorMessage name={`choices`} render={message => (
                                                            !(typeof message === 'string') && <div className={styles.error}>{message[idx]}</div>
                                                        )}></ErrorMessage>}
                                                    </div>
                                                )
                                            })}
                                            <ErrorMessage name={`choices`} render={message => (
                                                (typeof message === 'string') && <div className={styles.error}>{message}</div>
                                            )}></ErrorMessage>
                                            <button onClick={(event) => {
                                                event.preventDefault()
                                                push("")
                                            }}>add row</button>
                                        </div>
                                    )
                                }</FieldArray>
                                <div className={styles.options}>
                                    <h4>Options</h4>
                                    <div>
                                        <label htmlFor="allowMultiple">Multiple choices</label>
                                        <Field id="allowMultiple" as="select" name="settings.allowMultiple">
                                            {Array.from(Array(values.choices.length).keys(), (_, i) => i + 1).map(i => (
                                                <option value={i}>{i}</option>
                                            ))}
                                        </Field>
                                    </div>
                                    <div>
                                        <label htmlFor="timeout">Timeout in minutes</label>
                                        <Field id="timeout" name="settings.timeoutMinutes" type="number" min="1" max="120"></Field>
                                        <ErrorMessage name="settings.timeoutMinutes" render={message => (
                                            <div className={styles.error}>{message}</div>
                                        )}></ErrorMessage>
                                    </div>
                                    <div>
                                        <label htmlFor="filter">Filter</label>
                                        <Field id="filter" name="settings.filter" as="select">
                                            <option value="">none</option>
                                            <option value="ip">ip</option>
                                            <option value="cookie">cookie</option>
                                        </Field>
                                    </div>
                                </div>
                                <input type="submit"></input>
                            </Form>
                        )
                    }
                } validationSchema={validationSchema}
                ></Formik>
            </div>
        )
    }
}

export default CreatePoll