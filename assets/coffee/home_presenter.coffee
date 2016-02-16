class HomePresenter extends BasePresenter
  onCreate: () ->
    $('.js-input').on 'change', (event) =>
      @request.abort() if @request?

      keywords = $(event.target).val()
      query = keywords.replace(/\+/gi, "").replace(/(\s)+/gi, "+")

      @request = $.get
        url: "/search?keywords=" + query
        success: (data) =>
          console.log JSON.parse(data.posts)

