import React from 'react'
import { ApiGetPoll, ApiPutPollChoices } from '../../api/api'
import { RouteComponentProps } from 'react-router'
import { IPoll } from '../../models/poll'
import styles from './PollPage.module.css'
import { Formik, Form, Field, ErrorMessage } from 'formik'


interface MatchParams {
    uuid: string,
}

interface IState {
    poll: IPoll | null,
    selected: number,
}

interface IFormValues {
    selected: boolean[],
}

interface IFormErrors {
    selected?: string,
}

class PollPage extends React.Component<RouteComponentProps<MatchParams>, IState> {
    constructor(props: Readonly<RouteComponentProps<MatchParams, import("react-router").StaticContext, import("history").History.UnknownFacade>>) {
        super(props)
        this.state = {
            poll: null,
            selected: 0,
        }
    }

    getPoll = (once: boolean = false) => {
        ApiGetPoll(this.props.match.params.uuid).then(r => {
            let expires = new Date(r.data.expires)
            let poll: IPoll = r.data
            poll.expires = expires
            this.setState({ poll })

            if (once) return

            if (!isExpired(expires)) {
                setTimeout(() => {
                    this.getPoll()
                }, 1000)
            }
        })
    }

    handleFormChange = (selected: boolean[]) => {
        let selectedCount = count(selected, true)
        if (selectedCount !== this.state.selected) {
            this.setState({ selected: selectedCount })
        }
    }

    componentDidMount = () => {
        this.getPoll()
    }

    formValidate = (values: IFormValues): IFormErrors => {
        let res: IFormErrors = {}
        if (choicesRemained(this.state, count(values.selected, true)) < 0) res.selected = `only ${this.state.poll?.allowMultiple} choices is allowed`
        return res
    }

    onFormSubmit = ({ selected }: IFormValues) => {
        if (!this.state.poll) return
        let choices: string[] = []
        selected.forEach((choice, idx) => {
            if (!choice) return
            choices.push(this.state.poll!.choices[idx].text)
        })
        ApiPutPollChoices(this.state.poll.uuid, choices).then(() => { this.getPoll(true) }).catch(err => alert("couldn't commit votes: " + err))
    }

    render = () => {
        if (!this.state.poll) return null
        let choicesRem = choicesRemained(this.state)
        const pollExpired = isExpired(this.state.poll.expires)
        return (
            <div>
                <h2>{this.state.poll.name}</h2>
                {!pollExpired && <div className={"" + (choicesRem < 0 ? styles.red : "")}>Choices remained: {choicesRem}</div>}

                <Formik initialValues={{ selected: Array(this.state.poll.choices.length).fill(false) } as IFormValues} onSubmit={this.onFormSubmit}
                    validate={this.formValidate}
                >
                    {({ values }) => {
                        this.handleFormChange(values.selected)
                        return (
                            <Form>
                                <div>
                                    <table>
                                        {this.state.poll!.choices.map((choice, idx) => (
                                            <tr>
                                                <td>
                                                    <div key={choice.text} className={styles.inline}>{choice.text}:{choice.votes}</div>
                                                </td>
                                                <td>
                                                    {!pollExpired && <Field type="checkbox" name={`selected.${idx}`}></Field>}
                                                </td>
                                            </tr>
                                        ))}
                                    </table>
                                    <div className={styles.error}>
                                        <ErrorMessage name={"selected"}></ErrorMessage>
                                    </div>
                                    {!pollExpired && <div>
                                        <input type="submit"></input>
                                    </div>}
                                </div>
                            </Form>
                        )
                    }}
                </Formik>
                {pollExpired
                    ? <div>Poll is expired.</div>
                    : <div>{secondsRemained(this.state.poll.expires)} seconds remained</div>
                }
            </div>
        )
    }
}

const choicesRemained = (state: IState, selected?: number): number => {
    return (state.poll?.allowMultiple || 0) - (selected == null ? state.selected : selected)
}

const isExpired = (date: Date): boolean => {
    return date.valueOf() < Date.now()
}

const secondsRemained = (date: Date): number => {
    return Math.round((date.valueOf() - Date.now()) / 1000)
}

function count<T>(arr: T[], etalon: T): number {
    return arr.reduce((acc, v) => {
        if (v === etalon) acc++;
        return acc
    }, 0)
}

export default PollPage