import { FormLabel } from '@chakra-ui/form-control'
import { Input } from '@chakra-ui/react'

type EmailInputProps = {
  email: string
  changeFunc: (e: React.ChangeEvent<HTMLInputElement>, cred: string) => void
}

const EmailInput = ({ email, changeFunc }: EmailInputProps) => {
  return (
    <FormLabel width="100%">
      <Input
        placeholder="email"
        type="email"
        value={email}
        onChange={(e) => changeFunc(e, 'email')}
        required
        borderRadius="10px"
        borderColor="rgba(52, 72, 255, 1)"
      />
    </FormLabel>
  )
}

export default EmailInput
