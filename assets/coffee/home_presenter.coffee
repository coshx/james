class HomePresenter extends BasePresenter
  onCreate: () ->
    $('.js-input').on 'input change', (event) =>
      clearTimeout(@timer) if @timer?
      @request.abort() if @request?

      @timer = setTimeout(( () =>
        @timer = null
        @request.abort() if @request?

        keywords = $(event.target).val()
        query = keywords.replace(/\+/gi, "").replace(/(\s)+/gi, "+")

        @request = $.get
          url: "/search?keywords=" + query
          success: (data) =>
            @request = null
            console.log JSON.parse(data.posts)
    ), 300)

